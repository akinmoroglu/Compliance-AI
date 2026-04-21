echo "Running Test Case 1: Weight Loss"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Lose 20kg in 30 days with our supplement!", "headline": "Get the Supplement"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 2: Personal Attribute (Debt)"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Are you struggling with debt? We can help.", "headline": "Debt Help"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 3: Gambling"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Bet on Premier League today — big wins await!", "headline": "Play Now"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 4: Clean Ad"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Get fit this summer. Shop our activewear.", "headline": "Shop Now"}' http://localhost:8080/checks | jq .
