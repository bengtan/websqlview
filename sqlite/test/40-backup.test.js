runTest = async function() {
    const tempFilename = './foo.db'
    const db = new gosqlite.Database(':memory:')
    await db.open()
    await db.exec(`CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, word TEXT)`)
    await db.exec(`INSERT INTO foo(word) VALUES(?)`, 'Lorem')
    await db.exec(`INSERT INTO foo VALUES(?, ?)`, 10, 'ipsum')
    await db.exec(`INSERT INTO foo(word) VALUES(?)`, 'dolor')

    db2 = await db.backupTo(tempFilename)
    data = await db2.query(`SELECT id, word from foo`)
    if (data.length != 3) {
        return `Expected: 3 results, Actual: ${data.length}`
    }
    if (data[0].id != 1 || data[0].word != 'Lorem') {
        return `Unexpected data[0], Actual: ${JSON.stringify(data[0])}`
    }
    if (data[1].id != 10 || data[1].word != 'ipsum') {
        return `Unexpected data[1], Actual: ${JSON.stringify(data[1])}`
    }
    if (data[2].id != 11 || data[2].word != 'dolor') {
        return `Unexpected data[2], Actual: ${JSON.stringify(data[2])}`
    }

    await native.remove(tempFilename)
}
