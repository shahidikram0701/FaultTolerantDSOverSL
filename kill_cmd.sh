PID=$(ps aux | grep "./scalog.*--config" | head -1 | awk '{print $2}')
kill -9 $PID || true
