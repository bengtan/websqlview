gomysql = (function() {
    const Database = function(filename, id = -1) {
        this._id = id
        this.filename = filename
    }

    Database.prototype.open = function() {
        return new Promise((resolve, reject) => {
            _mysqlMux('open', this.filename).then(id => {
                this._id = id
                resolve()
            }).catch(reject)
        })
    }

    Database.prototype.close = function() {
        return new Promise((resolve, reject) => {
            const handle = this._id
            this._id = -1
            _mysqlMux('close', handle)
                .then(resolve).catch(reject)
        })
    }

    Database.prototype.exec = function(query, ...params) {
        return _mysqlMux('exec', this._id, query, ...params)
    }
    Database.prototype.query = function(query, ...params) {
        return _mysqlMux('query', this._id, query, ...params)
    }
    Database.prototype.queryRow = function(query, ...params) {
        return _mysqlMux('queryRow', this._id, query, ...params)
    }
    Database.prototype.queryResult = function(query, ...params) {
        return _mysqlMux('queryResult', this._id, query, ...params)
    }

    Database.prototype.begin = function() {
        return new Promise((resolve, reject) => {
            _mysqlMux('begin', this._id).then(id => {
                tx = new Transaction(id)
                resolve(tx)
            }).catch(reject)
        })
    }

    const Transaction = function(id = -1) {
        this._id = id
    }

    Transaction.prototype.commit = function() {
        return new Promise((resolve, reject) => {
            const handle = this._id
            this._id = -1
            _mysqlMux('tx.commit', handle)
                .then(resolve).catch(reject)
        })
    }

    Transaction.prototype.rollback = function() {
        return new Promise((resolve, reject) => {
            const handle = this._id
            this._id = -1
            _mysqlMux('tx.rollback', handle)
                .then(resolve).catch(reject)
        })
    }

    Transaction.prototype.exec = function(query, ...params) {
        return _mysqlMux('tx.exec', this._id, query, ...params)
    }
    Transaction.prototype.query = function(query, ...params) {
        return _mysqlMux('tx.query', this._id, query, ...params)
    }
    Transaction.prototype.queryRow = function(query, ...params) {
        return _mysqlMux('tx.queryRow', this._id, query, ...params)
    }
    Transaction.prototype.queryResult = function(query, ...params) {
        return _mysqlMux('tx.queryResult', this._id, query, ...params)
    }

    return {
        Database
    }
})()
