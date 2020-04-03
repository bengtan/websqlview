runTest = async function() {
    const db1 = new sqlite.Database(':memory:')
    await db1.open()
    await db1.close()

    const db2 = new sqlite.Database(':memory:')
    await db2.open()
    if (db2._id != 1) {
        return `Expected: 1, Actual: ${db2._id}`
    }
    await db2.close()

    const db3 = new sqlite.Database('/foo.db')
    try {
        await db3.open()
        return `Opened database in root directory`
    } catch (e) {}
}
