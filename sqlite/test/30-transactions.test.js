runTest = async function() {
    const db = new gosqlite.Database(':memory:')
    await db.open()
    await db.exec(`CREATE TABLE foo (id INTEGER NOT NULL PRIMARY KEY, word TEXT)`)
    await db.exec(`INSERT INTO foo(word) VALUES(?)`, 'Lorem')

    // Test: Rollback of a simple insert
    tx1 = await db.begin()
    handle1 = tx1._id
    await tx1.exec(`INSERT INTO foo(word) VALUES(?)`, 'ipsum')
    data = await tx1.queryResult(`SELECT COUNT(*) FROM foo`)
    if (data != 2) {
        return `Expected: COUNT(*) == 2, Actual: ${data}`
    }

    await tx1.rollback()
    data = await db.queryResult(`SELECT COUNT(*) FROM foo`)
    if (data != 1) {
        return `Expected: COUNT(*) == 1, Actual: ${data}`
    }

    // Test: Commit of a simple insert
    tx2 = await db.begin()
    handle2 = tx2._id
    await tx2.exec(`INSERT INTO foo(word) VALUES(?)`, 'ipsum')
    await tx2.commit()

    data = await db.queryResult(`SELECT COUNT(*) FROM foo`)
    if (data != 2) {
        return `Expected: COUNT(*) == 2, Actual: ${data}`
    }

    // Test: Expect the handle to be reused
    if (handle1 != handle2) {
        return `Expected: ${handle1}, Actual: ${handle2}`
    }

    await db.close()
}
