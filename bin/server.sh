echo "Starting the Server"
if ! command -v air &> /dev/null; then
    echo "Error: 'air' is not installed. Installing it now..."
    go install github.com/air-verse/air
fi

air