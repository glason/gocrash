rm android.json
echo getting data on 10.64.12.213
curl http://10.64.12.213:8080/logs/clientCrash.json >> android.json
echo getting data on 10.64.13.226
curl http://10.64.13.226:8080/logs/clientCrash.json >> android.json
echo getting data on 10.64.11.188
curl http://10.64.11.188:8080/logs/clientCrash.json >> android.json
echo getting data on 10.64.11.187
curl http://10.64.11.187:8080/logs/clientCrash.json >> android.json
grep '"ct":"android' android.json > crash.json
