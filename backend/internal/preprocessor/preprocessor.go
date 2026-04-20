package preprocessor

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"

    "compliance-checker/internal/models"
)

const maxVideoSizeMB = 200

// Output holds everything the pre-processor extracted from the HTTP request.
// This is the input to all Phase 1 goroutines.
type Output struct {
    Req         *models.CheckRequest
    ImageBytes  []byte   // nil if no image
    ImageSize   int64
    VideoFrames [][]byte // nil if no video; JPEG bytes per frame
    AudioBytes  []byte   // nil if no video or no audio track found
    VideoSize   int64
    HasImage    bool
    HasVideo    bool
    VideoDurationSeconds float64 // 0 if not a video
}

// Parse reads the HTTP request and extracts all inputs.
// Handles both multipart/form-data (when a file is present) and
// application/json (text-only, backward-compatible).
func Parse(r *http.Request) (*Output, error) {
    out := &Output{Req: &models.CheckRequest{}}

    contentType := r.Header.Get("Content-Type")

    if strings.Contains(contentType, "application/json") {
        // JSON path — text-only, no file
        if err := parseJSON(r, out.Req); err != nil {
            return nil, fmt.Errorf("json parse: %w", err)
        }
        return out, nil
    }

    // Multipart path
    if err := r.ParseMultipartForm(maxVideoSizeMB * 1024 * 1024); err != nil {
        return nil, fmt.Errorf("multipart parse: %w", err)
    }
    if err := parseFormFields(r, out.Req); err != nil {
        return nil, fmt.Errorf("form fields: %w", err)
    }

    // Check for image file
    imageFile, imageHeader, err := r.FormFile("image")
    if err == nil {
        defer imageFile.Close()
        out.ImageBytes, err = io.ReadAll(imageFile)
        if err != nil {
            return nil, fmt.Errorf("read image: %w", err)
        }
        out.ImageSize = imageHeader.Size
        out.HasImage = true
    }

    // Check for video file
    videoFile, videoHeader, err := r.FormFile("video")
    if err == nil {
        defer videoFile.Close()

        if videoHeader.Size > maxVideoSizeMB*1024*1024 {
            return nil, fmt.Errorf("video exceeds %dMB limit", maxVideoSizeMB)
        }

        videoBytes, err := io.ReadAll(videoFile)
        if err != nil {
            return nil, fmt.Errorf("read video: %w", err)
        }

        // Detect file extension from filename for ffmpeg
        ext := filepath.Ext(videoHeader.Filename)
        if ext == "" {
            ext = ".mp4" // safe default
        }

        frames, audioBytes, duration, err := extractVideoContent(videoBytes, ext)
        if err != nil {
            return nil, fmt.Errorf("video extraction: %w", err)
        }

        out.VideoFrames = frames
        out.AudioBytes = audioBytes
        out.VideoSize = videoHeader.Size
        out.VideoDurationSeconds = duration
        out.HasVideo = true
    }

    return out, nil
}

// parseFormFields reads all text fields from the multipart form into CheckRequest.
func parseFormFields(r *http.Request, req *models.CheckRequest) error {
    req.Platform = r.FormValue("platform")
    req.Region = r.FormValue("region")
    req.PrimaryText = r.FormValue("primary_text")
    req.Headline = r.FormValue("headline")
    req.Description = r.FormValue("description")
    req.LandingPageURL = r.FormValue("landing_page_url")
    req.ForceRefresh = r.FormValue("force_refresh") == "true"
    req.AdFormat = r.FormValue("ad_format")

    ageMin, _ := strconv.Atoi(r.FormValue("age_min"))
    ageMax, _ := strconv.Atoi(r.FormValue("age_max"))
    req.AgeMin = ageMin
    req.AgeMax = ageMax

    // Countries come as a JSON array string: '["US","GB","TR"]'
    // Parse it; on failure treat as global (empty slice)
    countriesRaw := r.FormValue("countries")
    if countriesRaw != "" && countriesRaw != "[]" {
        // Strip brackets and split by comma — simple parse, avoids json import cycle
        trimmed := strings.Trim(countriesRaw, "[]")
        for _, c := range strings.Split(trimmed, ",") {
            c = strings.Trim(strings.TrimSpace(c), `"`)
            if c != "" {
                req.Countries = append(req.Countries, c)
            }
        }
    }

    return nil
}

// parseJSON reads a JSON body into CheckRequest.
func parseJSON(r *http.Request, req *models.CheckRequest) error {
    return json.NewDecoder(r.Body).Decode(req)
}

// extractVideoContent runs ffmpeg to extract frames and audio from a video file.
func extractVideoContent(videoBytes []byte, ext string) (frames [][]byte, audioBytes []byte, duration float64, err error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, nil, 0, fmt.Errorf("ffmpeg not found in PATH")
	}

    tmpDir, err := os.MkdirTemp("", "compliance-video-*")
    if err != nil {
        return nil, nil, 0, err
    }
    defer os.RemoveAll(tmpDir)

    // Write video bytes to temp file
    videoPath := filepath.Join(tmpDir, "input"+ext)
    if err := os.WriteFile(videoPath, videoBytes, 0644); err != nil {
        return nil, nil, 0, err
    }

    // Get video duration via ffprobe
    duration = getVideoDuration(videoPath)

    // Step 1: Scene-change frame extraction
    framesDir := filepath.Join(tmpDir, "frames")
    os.MkdirAll(framesDir, 0755)

    cmd := exec.Command("ffmpeg",
        "-i", videoPath,
        "-vf", "select='gt(scene,0.4)'",
        "-vsync", "vfr",
        "-q:v", "2",
        "-frames:v", "30", // hard cap
        filepath.Join(framesDir, "frame_%03d.jpg"),
    )
    cmd.Run() // error is non-fatal — check frame count below

    frameFiles, _ := filepath.Glob(filepath.Join(framesDir, "*.jpg"))

    // Step 2: Fallback to uniform 1fps if fewer than 3 frames
    if len(frameFiles) < 3 {
        os.RemoveAll(framesDir)
        os.MkdirAll(framesDir, 0755)

        maxFrames := int(duration) // 1fps for duration gives N frames
        if maxFrames > 30 {
            maxFrames = 30
        }
        if maxFrames < 1 {
            maxFrames = 10 // minimum attempt
        }

        cmd = exec.Command("ffmpeg",
            "-i", videoPath,
            "-r", "1",
            "-q:v", "2",
            "-frames:v", strconv.Itoa(maxFrames),
            filepath.Join(framesDir, "frame_%03d.jpg"),
        )
        cmd.Run()
        frameFiles, _ = filepath.Glob(filepath.Join(framesDir, "*.jpg"))
    }

    // Cap at 30 frames
    if len(frameFiles) > 30 {
        frameFiles = frameFiles[:30]
    }

    for _, f := range frameFiles {
        b, readErr := os.ReadFile(f)
        if readErr != nil {
            continue
        }
        frames = append(frames, b)
    }

    // Extract audio track as MP3
    audioPath := filepath.Join(tmpDir, "audio.mp3")
    audioCmd := exec.Command("ffmpeg",
        "-i", videoPath,
        "-q:a", "0",
        "-map", "a",
        audioPath,
    )
    if audioErr := audioCmd.Run(); audioErr == nil {
        audioBytes, _ = os.ReadFile(audioPath)
    }
    // Audio extraction failure is non-fatal (video may have no audio track)

    return frames, audioBytes, duration, nil
}

// getVideoDuration returns the duration in seconds using ffprobe.
// Returns 0 on failure — callers treat 0 as unknown.
func getVideoDuration(videoPath string) float64 {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return 0
	}
    cmd := exec.Command("ffprobe",
        "-v", "error",
        "-show_entries", "format=duration",
        "-of", "csv=p=0",
        videoPath,
    )
    out, err := cmd.Output()
    if err != nil {
        return 0
    }
    d, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
    return d
}
