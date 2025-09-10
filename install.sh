#!/bin/bash
set -e

INSTALL_DIR=${1:-/usr/local/bin}
BINARY_NAME=${TODO_NAME:-todo}


echo "Building the application..."
go build -o "$BINARY_NAME"

echo "Copying binary to '$INSTALL_DIR' (may require sudo)..."

if [ ! -d "$INSTALL_DIR" ]; then
    echo "Directory '$INSTALL_DIR' does not exist. Creating..."
    if [ "$(id -u)" -ne 0 ]; then
        sudo mkdir -p "$$INSTALL_DIR"
    else
        mkdir -p "$$INSTALL_DIR"
    fi
fi

if [ -w "$INSTALL_DIR" ]; then
    cp "$BINARY_NAME" "$INSTALL_DIR/"
else
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
fi

if [ -w "$INSTALL_DIR" ]; then
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "Installation completed successfully."

# Проверяем, есть ли каталог установки в $PATH и выводим подсказку, если нет
if ! echo "$PATH" | grep -Eq "(^|:)$INSTALL_DIR($|:)"; then
    echo "Note: '$INSTALL_DIR' is not in your PATH."
    echo "To run '$BINARY_NAME' command easily, add the directory to your PATH:"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    echo "Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc)."
fi

echo "You can now run the application by typing '$BINARY_NAME' anywhere."