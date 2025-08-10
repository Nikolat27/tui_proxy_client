#!/bin/bash

# Create a simple 256x256 PNG icon using ImageMagick
# If ImageMagick is not available, we'll create a text-based placeholder

if command -v convert >/dev/null 2>&1; then
    # Create a simple icon with text
    convert -size 256x256 xc:transparent \
        -fill "#4A90E2" \
        -draw "rectangle 20,20 236,236" \
        -fill white \
        -pointsize 48 \
        -gravity center \
        -draw "text 0,0 'TUI\nProxy'" \
        tui_proxy_client.png
else
    echo "ImageMagick not found. Creating a placeholder icon..."
    # Create a simple colored square as placeholder
    echo "iVBORw0KGgoAAAANSUhEUgAAAQAAAAEACAYAAABccqhmAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAALEwAACxMBAJqcGAAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3Njape.org5vuPBoAAA==" | base64 -d > tui_proxy_client.png
fi 