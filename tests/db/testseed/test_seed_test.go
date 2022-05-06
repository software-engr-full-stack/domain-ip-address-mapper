package testseed

import (
    "testing"

    "demo/lib/testlib"
    "demo/lib/testlib/common"
    "demo/lib/testlib/common/reddit"

    "demo/db/testseed"
    "demo/entities/crobat"
    "demo/entities/user"
)

func TestSeedEndToEnd(t *testing.T) {
    t.Parallel()
    setup := testlib.Setup()
    setup.TestDB.Setup()
    defer setup.TestDB.CleanUp()

    instance, err := setup.DBObj.Open()
    if err != nil {
        panic(err)
    }
    defer func() {
        if err = setup.DBObj.Close(); err != nil {
            panic(err)
        }
    }()

    ctx := setup.ContextWithDBInstance(instance)

    in := crobat.InputType{
        DomainSub:  "reddit.com",
        UniqueSort: true,
    }
    err = testseed.TestSeed(ctx, in)
    if err != nil {
        panic(err)
    }

    users, err := user.All(ctx)
    if err != nil {
        panic(err)
    }
    common.TestUsers(t, ctx, users)

    rootUser, isFound, err := user.Root(ctx)
    if err != nil {
        panic(err)
    }
    common.TestRootUser(t, ctx, rootUser, isFound)

    reddit.TestData(t, ctx, rootUser.Data)
    reddit.TestSubDomains(t, ctx, rootUser.Data)
}
