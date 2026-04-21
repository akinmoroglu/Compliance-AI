echo "Running Test Case 5: Health Claim (Doctors Hate Him)"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "DOCTORS HATE HIM! This one trick cures arthritis!", "headline": "Heal Fast"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 6: Impersonation (Rolex)"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Buy Rolex watches — genuine Swiss luxury", "headline": "Luxury Watches"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 7: Urgency Tactic"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Only 3 spots left! Sign up today or miss out!", "headline": "Last Chance"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 8: Crypto Exchange"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Trade Bitcoin on our platform — high returns guaranteed", "headline": "Crypto Wealth"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 9: Body Shaming"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "Tired of your belly fat? Our pill melts it away in 7 days!", "headline": "Melts Fat"}' http://localhost:8080/checks | jq .

echo -e "\nRunning Test Case 10: Christian targeting"
curl -s -X POST -H "Content-Type: application/json" -d '{"primary_text": "For Christians looking to invest ethically.", "headline": "Ethical Investing"}' http://localhost:8080/checks | jq .
