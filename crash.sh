rm android.json
touch crash_old.json
echo getting data on 10.64.12.213
curl http://10.64.12.213:8080/logs/clientCrash.json >> android.json
echo getting data on 10.64.13.226
curl http://10.64.13.226:8080/logs/clientCrash.json >> android.json
echo getting data on 10.64.11.188
curl http://10.64.11.188:8080/logs/clientCrash.json >> android.json
echo getting data on 10.64.11.187
curl http://10.64.11.187:8080/logs/clientCrash.json >> android.json
grep '"ct":"android' android.json > crash_new.json
diff -w crash_old.json crash_new.json | grep '>  INFO' > crash.json
mv crash_new.json crash_old.json