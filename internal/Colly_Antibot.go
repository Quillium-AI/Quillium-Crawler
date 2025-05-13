c.Limit(&colly.LimitRule{
    DomainGlob:  "*",
    Delay:       2 * time.Second,
    RandomDelay: 1 * time.Second,
})