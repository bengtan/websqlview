runTest = async function() {
    const db1 = new sqlite.Database(':memory:')
    await db1.open()
    const handle1 = db1._id
    await db1.close()

    const db2 = new sqlite.Database(':memory:')
    await db2.open()

    // Expect the handle to be reused
    if (db2._id != handle1) {
        return `Expected: ${handle1}, Actual: ${db2._id}`
    }
    await db2.close()

    // This should fail
    const db3 = new sqlite.Database('/foo.db')
    try {
        await db3.open()
        return `Unexpected: Opened database in root directory`
    } catch (e) {}
}
