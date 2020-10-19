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
        writeFile: (filename, blob) => {
            return new Promise((resolve, reject) => {
                const reader = new FileReader()
                reader.addEventListener('loadend', () => {
                    if (reader.result) {
                        _nativeMux('writeFile', filename, reader.result.replace(/data:.*base64,/, ''))
                        .then(resolve)
                        .catch(reject)
                    }
                    else {
                        reject('readAsDataURL: No result')
                    }
                })
                reader.readAsDataURL(new Blob([blob], {type: 'application/octet-stream'}))
            })
        },
    }
})()
