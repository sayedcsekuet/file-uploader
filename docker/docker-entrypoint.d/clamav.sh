#!/usr/bin/env bash
function updateClamAvDatabase() {
  # Update files if they are missing or older that X days
  if [[ $(find "$CLAM_DATA_PATH/daily.cld" -mtime +2 -print) ]] || [[ ! -f "$CLAM_DATA_PATH/daily.cld" ]]; then
    echo "Clamd files are too old. Updating..."
    echo "OK: Running freshclam to update virus databases. This can take a few minutes..."
    sleep 1
    freshclam
  else
    echo "File clamd files are current"
  fi

}

echo "Updating claman database"
updateClamAvDatabase

echo "Starting clamd service"
clamd
