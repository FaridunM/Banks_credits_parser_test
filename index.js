const { createServer } = require('@vercel/node')
const { join } = require('path')
const { spawn } = require('child_process')

const server = createServer((req, res) => {
  const child = spawn('./main', [], { cwd: join(__dirname, 'api') })
  child.stdout.pipe(res)
})

module.exports = server