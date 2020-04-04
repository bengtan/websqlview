runTest = async function() {
    const db = new sqlite.Database(':memory:')
    await db.open()
    await db.exec(`CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, word TEXT)`)
    
    data = await db.exec(`INSERT INTO foo(word) VALUES(?)`, 'Lorem')
    if (data.lastInsertId != 1) {
        return `Expected: lastInsertId == 1, Actual: ${data.lastInsertId}`
    }
    data = await db.exec(`INSERT INTO foo VALUES(?, ?)`, 10, 'ipsum')
    if (data.lastInsertId != 10) {
        return `Expected: lastInsertId == 10, Actual: ${data.lastInsertId}`
    }
    data = await db.exec(`INSERT INTO foo(word) VALUES(?)`, 'dolor')
    if (data.lastInsertId != 11) {
        return `Expected: lastInsertId == 11, Actual: ${data.lastInsertId}`
    }

    data = await db.query(`SELECT id, word from foo`)
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

    data = await db.queryRow(`SELECT id, word from foo`)
    if (data.id != 1 || data.word != 'Lorem') {
        return `Unexpected data, Actual: ${JSON.stringify(data)}`
    }

    data = await db.queryResult(`SELECT word from foo WHERE id = ?`, 1)
    if (data != 'Lorem') {
        return `Expected Lorem, Actual: ${data}`
    }

    // This should fail
    try {
        await db.queryResult(`SELECT id, word from foo WHERE id = ?`, 1)
        return `Unexpected: queryResult() returned`
    } catch (e) {}

    await db.exec(`UPDATE foo SET word = ? WHERE id = ?`, 'LOREM', 1)
    data = await db.queryResult(`SELECT word from foo WHERE id = ?`, 1)
    if (data != 'LOREM') {
        return `Expected LOREM, Actual: ${data}`
    }

    await db.exec(`DELETE FROM foo WHERE id = ?`, 1)
    data = await db.queryResult(`SELECT COUNT(id) from foo`)
    if (data != 2) {
        return `Expected: COUNT(id) == 2, Actual: ${data}`
    }

    try {
        await db.queryResult(`SELECT word from foo WHERE id = ?`, 1)
        return `Unexpected: queryResult() returned a value`
    } catch (e) {}

    // ToDo: upsert?
    await db.close()
}
