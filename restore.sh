#!/bin/bash
# This script restores dotfiles and installs Homebrew packages.
# Exit immediately if a command exits with a non-zero status.
set -e

# Find the script's own directory to reliably locate other files.
DOTFILES_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CONFIG_SOURCE_DIR="$DOTFILES_DIR/config"
CONFIG_DEST_DIR="$HOME/.config"

# --- Homebrew Package Installation ---
echo "--- Checking for Homebrew ---"
if ! command -v brew &> /dev/null; then
    echo "Warning: Homebrew not found. Skipping package installation."
    echo "Please install Homebrew first by visiting https://brew.sh/"
else
    echo "--- Installing Homebrew Formulae ---"
    FORMULAE_FILE="$DOTFILES_DIR/brew_formulae.txt"
    if [ -f "$FORMULAE_FILE" ]; then
        # Read file line by line and install
        while read -r formula || [ -n "$formula" ]; do
            # Process only non-empty lines
            if [ -n "$formula" ]; then
                echo "Installing formula: $formula"
                # Use '|| true' to continue even if a formula is already installed or fails
                brew install "$formula" || echo "Warning: Failed to install '$formula'. It may already be installed."
            fi
        done < "$FORMULAE_FILE"
    else
        echo "Warning: brew_formulae.txt not found. Skipping formulae."
    fi

    echo ""
    echo "--- Installing Homebrew Casks ---"
    CASKS_FILE="$DOTFILES_DIR/brew_casks.txt"
    if [ -f "$CASKS_FILE" ]; then
        while read -r cask || [ -n "$cask" ]; do
            if [ -n "$cask" ]; then
                echo "Installing cask: $cask"
                brew install --cask "$cask" || echo "Warning: Failed to install '$cask'. It may already be installed."
            fi
        done < "$CASKS_FILE"
    else
        echo "Warning: brew_casks.txt not found. Skipping casks."
    fi
fi

# --- Dotfile Symlinking ---
echo ""
echo "--- Starting Dotfile Symlinking ---"
mkdir -p "$CONFIG_DEST_DIR" # Ensure ~/.config directory exists

# Loop through all items in the source config directory
for source_path in "$CONFIG_SOURCE_DIR"/*; do
    # Only process directories
    if [ -d "$source_path" ]; then
        config_name=$(basename "$source_path")
        dest_path="$CONFIG_DEST_DIR/$config_name"

        echo "Processing $config_name..."

        # If a link or directory already exists at the destination, remove it.
        if [ -e "$dest_path" ] || [ -L "$dest_path" ]; then
            echo "  - Removing existing target: $dest_path"
            rm -rf "$dest_path"
        fi

        # Create the symbolic link
        echo "  - Linking $source_path -> $dest_path"
        ln -s "$source_path" "$dest_path"
        echo "  - Successfully linked $config_name."
    fi
done

echo ""
echo "--- Restore Script Complete ---"
