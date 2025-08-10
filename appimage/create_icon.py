#!/usr/bin/env python3

try:
    from PIL import Image, ImageDraw, ImageFont
    import os
    
    # Create a 256x256 image with a blue background
    img = Image.new('RGBA', (256, 256), (74, 144, 226, 255))
    draw = ImageDraw.Draw(img)
    
    # Draw a terminal-like rectangle
    draw.rectangle([48, 68, 208, 188], fill=(44, 62, 80, 255), outline=(52, 73, 94, 255), width=2)
    
    # Draw terminal header
    draw.rectangle([48, 68, 208, 88], fill=(52, 73, 94, 255))
    
    # Draw terminal buttons (red, yellow, green)
    draw.ellipse([64, 74, 72, 82], fill=(231, 76, 60, 255))  # Red
    draw.ellipse([84, 74, 92, 82], fill=(243, 156, 18, 255))  # Yellow
    draw.ellipse([104, 74, 112, 82], fill=(39, 174, 96, 255))  # Green
    
    # Draw terminal text lines
    colors = [(236, 240, 241, 255)] * 7
    widths = [80, 60, 90, 70, 50, 100, 40]
    
    for i, (color, width) in enumerate(zip(colors, widths)):
        y = 100 + i * 10
        draw.rounded_rectangle([58, y, 58 + width, y + 4], radius=2, fill=color)
    
    # Draw cursor
    draw.rounded_rectangle([58, 170, 66, 174], radius=2, fill=(231, 76, 60, 255))
    
    # Draw network nodes
    draw.ellipse([168, 88, 192, 112], fill=(52, 152, 219, 200))  # Blue node
    draw.ellipse([192, 112, 208, 128], fill=(155, 89, 182, 200))  # Purple node
    draw.ellipse([182, 132, 202, 152], fill=(230, 126, 34, 200))  # Orange node
    
    # Draw connection lines
    draw.line([180, 100, 200, 120], fill=(149, 165, 166, 150), width=2)
    draw.line([200, 120, 190, 140], fill=(149, 165, 166, 150), width=2)
    draw.line([190, 140, 180, 100], fill=(149, 165, 166, 150), width=2)
    
    # Save the image
    img.save('tui_proxy_client.png')
    print("Icon created successfully: tui_proxy_client.png")
    
except ImportError:
    print("PIL/Pillow not available. Creating a simple colored square...")
    # Fallback: create a simple colored square
    with open('tui_proxy_client.png', 'wb') as f:
        # Simple 1x1 blue PNG
        png_data = b'\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\x0cIDATx\x9cc```\x00\x00\x00\x04\x00\x01\xf5\x01\x00\x00\x00\x00IEND\xaeB`\x82'
        f.write(png_data)
    print("Simple fallback icon created") 