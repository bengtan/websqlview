// Note: This assumes that a pgsql database 'test' exists, and is empty.
runTest = async function() {
    const db = new gopgsql.Database(DB_DSN)
    await db.open()
    await db.exec(`DROP TABLE IF EXISTS foo`)
    await db.exec(`CREATE TABLE foo (id serial, word text)`)

    data = await db.query(`INSERT INTO foo (word) VALUES ($1) RETURNING id`, 'Lorem')
    if (data.length != 1 || data[0].id != 1) {
        return `Expected: data[0].id == 1, Actual data: ${JSON.stringify(data)}`
    }

    data = await db.exec(`INSERT INTO foo (word) VALUES ($1), ($2)`, 'ipsum', 'dolor')
    data = await db.query(`INSERT INTO foo (word) VALUES ($1), ($2) RETURNING *`, 'sit', 'amet')
    if (data.length != 2 || data[0].id != 4 || data[1].id != 5) {
        return `Expected: data[0].id == 4 and data[1] == 5, Actual data: ${JSON.stringify(data)}`
    }

    data = await db.query(`SELECT id, word from foo`)
    if (data.length != 5) {
        return `Expected: 5 results, Actual: ${data.length}`
    }
    if (data[0].id != 1 || data[0].word != 'Lorem') {
        return `Unexpected data[0], Actual: ${JSON.stringify(data[0])}`
    }
    if (data[1].id != 2 || data[1].word != 'ipsum') {
        return `Unexpected data[1], Actual: ${JSON.stringify(data[1])}`
    }
    if (data[2].id != 3 || data[2].word != 'dolor') {
        return `Unexpected data[2], Actual: ${JSON.stringify(data[2])}`
    }
    if (data[3].id != 4 || data[3].word != 'sit') {
        return `Unexpected data[3], Actual: ${JSON.stringify(data[3])}`
    }
    if (data[4].id != 5 || data[4].word != 'amet') {
        return `Unexpected data[4], Actual: ${JSON.stringify(data[4])}`
    }

    data = await db.queryRow(`SELECT id, word from foo`)
    if (data.id != 1 || data.word != 'Lorem') {
        return `Unexpected data, Actual: ${JSON.stringify(data)}`
    }

    data = await db.queryRow(`SELECT id, word from foo WHERE FALSE`)
    if (data) {
        return `Expected NULL, Actual: ${JSON.stringify(data)}`
    }

    data = await db.queryResult(`SELECT word from foo WHERE id = $1`, 1)
    if (data != 'Lorem') {
        return `Expected Lorem, Actual: ${data}`
    }

    // This should fail because queryResult() only returns one value
    try {
        await db.queryResult(`SELECT id, word from foo WHERE id = $1`, 1)
        return `Unexpected: queryResult() returned`
    } catch (e) {}

    await db.exec(`UPDATE foo SET word = $1 WHERE id = $2`, 'LOREM', 1)
    data = await db.queryResult(`SELECT word from foo WHERE id = $1`, 1)
    if (data != 'LOREM') {
        return `Expected LOREM, Actual: ${data}`
    }

    await db.exec(`DELETE FROM foo WHERE id = $1`, 1)
    data = await db.queryResult(`SELECT COUNT(id) from foo`)
    if (data != 4) {
        return `Expected: COUNT(id) == 4, Actual: ${data}`
    }

    try {
        await db.queryResult(`SELECT word from foo WHERE id = ?`, 1)
        return `Unexpected: queryResult() returned a value`
    } catch (e) {}

    // ToDo: upsert?
    await db.close()
}
