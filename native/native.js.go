package native

const _nativeJs = `
native = (function() {
    return {
        setTitle: exitCode => {
            return _nativeMux('setTitle', exitCode)
        },
        exit: exitCode => {
            return _nativeMux('exit', exitCode)
        },
        remove: name => {
            return _nativeMux('remove', name)
        },
    }
})()
`
