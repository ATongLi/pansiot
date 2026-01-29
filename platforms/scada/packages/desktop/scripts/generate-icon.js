/**
 * 生成应用图标
 * 将 SVG Logo 转换为 PNG 格式
 */

const sharp = require('sharp')
const fs = require('fs')
const path = require('path')

const svgPath = path.join(__dirname, '../assets/icon.svg')
const pngPath = path.join(__dirname, '../assets/icon.png')
const sizes = [16, 32, 48, 64, 128, 256]

async function generateIcons() {
  try {
    // 读取 SVG 文件
    const svgBuffer = fs.readFileSync(svgPath)

    // 生成不同尺寸的 PNG
    for (const size of sizes) {
      const outputPath = path.join(__dirname, `../assets/icon-${size}x${size}.png`)

      await sharp(svgBuffer)
        .resize(size, size)
        .png()
        .toFile(outputPath)

      console.log(`✓ Generated ${path.basename(outputPath)}`)
    }

    // 生成默认 256x256 图标
    await sharp(svgBuffer)
      .resize(256, 256)
      .png()
      .toFile(pngPath)

    console.log(`✓ Generated ${path.basename(pngPath)}`)
    console.log('\n✅ All icons generated successfully!')
  } catch (error) {
    console.error('❌ Error generating icons:', error.message)
    process.exit(1)
  }
}

generateIcons()
