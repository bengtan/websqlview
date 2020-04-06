package sqlite

const _sqliteJs=`
sqlite = (function() {
    const Database = function(filename) {
        this._id = -1
        this.filename = filename
    }

    Database.prototype.open = function() {
        return new Promise((resolve, reject) => {
            _sqliteMux('open', this.filename).then((id) => {
                this._id = id
                resolve()
            }).catch(reject)
        })
    }

    Database.prototype.close = function() {
        return new Promise((resolve, reject) => {
            const handle = this._id
            this._id = -1
            _sqliteMux('close', handle)
                .then(resolve).catch(reject)
        })
    }

    Database.prototype.exec = function(query, ...params) {
        return _sqliteMux('exec', this._id, query, ...params)
    }
    Database.prototype.query = function(query, ...params) {
        return _sqliteMux('query', this._id, query, ...params)
    }
    Database.prototype.queryRow = function(query, ...params) {
        return _sqliteMux('queryRow', this._id, query, ...params)
    }
    Database.prototype.queryResult = function(query, ...params) {
        return _sqliteMux('queryResult', this._id, query, ...params)
    }

    Database.prototype.begin = function() {
        return new Promise((resolve, reject) => {
            _sqliteMux('begin', this._id).then((id) => {
                tx = new Transaction()
                tx._id = id
                resolve(tx)
            }).catch(reject)
        })
    }

    const Transaction = function() {
        this._id = -1
    }

    Transaction.prototype.commit = function() {
        return new Promise((resolve, reject) => {
            const handle = this._id
            this._id = -1
            _sqliteMux('tx.commit', handle)
                .then(resolve).catch(reject)
        })
    }

    Transaction.prototype.rollback = function() {
        return new Promise((resolve, reject) => {
            const handle = this._id
            this._id = -1
            _sqliteMux('tx.rollback', handle)
                .then(resolve).catch(reject)
        })
    }

    Transaction.prototype.exec = function(query, ...params) {
        return _sqliteMux('tx.exec', this._id, query, ...params)
    }
    Transaction.prototype.query = function(query, ...params) {
        return _sqliteMux('tx.query', this._id, query, ...params)
    }
    Transaction.prototype.queryRow = function(query, ...params) {
        return _sqliteMux('tx.queryRow', this._id, query, ...params)
    }
    Transaction.prototype.queryResult = function(query, ...params) {
        return _sqliteMux('tx.queryResult', this._id, query, ...params)
    }

    return {
        Database
    }
})()
`