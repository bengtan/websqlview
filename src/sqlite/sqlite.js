sqlite = {}

sqlite.Database = function Database(filename) {
    this._id = -1
    this.filename = filename
}

sqlite.Database.prototype.open = function() {
    return new Promise((resolve, reject) => {
        _sqliteMux('open', this.filename).then((id) => {
            this._id = id
            resolve()
        }).catch(reject)
    })
}

sqlite.Database.prototype.close = function() {
    return new Promise((resolve, reject) => {
        const handle = this._id
        this._id = -1
        _sqliteMux('close', handle)
            .then(resolve).catch(reject)
    })
}

sqlite.Database.prototype.exec = function(query, ...params) {
    return new Promise((resolve, reject) => {
        _sqliteMux('exec', this._id, query, ...params)
            .then((data) => {
                resolve({
                    lastInsertId: data[0],
                    rowsAffected: data[1],
                })
            }).catch(reject)
    })
}

