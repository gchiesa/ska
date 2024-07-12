// commitlint.configuration.js
module.exports = {
  extends: ['@commitlint/configuration-conventional'],
  ignores: [(message) => /^Bumps \[.+]\(.+\) from .+ to .+\.$/m.test(message)],
}
