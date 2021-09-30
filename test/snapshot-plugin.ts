module.exports = {
  test: (val: any) => typeof val === 'string',
  serialize: (val: any) => {
    return `"${val.replace(
      /AssetParameters([A-Fa-f0-9]{64})(\w+)|(\w+) (\w+) for asset\s?(version)?\s?"([A-Fa-f0-9]{64})"/,
      '[HASH REMOVED]',
    )}"`
  },
}
