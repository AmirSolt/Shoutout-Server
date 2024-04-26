Extensions:

- Added customer collection
- Added stripe webhook
- env.go and glitchtip.go are to manage the above extensions

Stripe notes:
- Create a webhook (see to stripewebhooks.go for the list of events)
- Create a Product/s with Price/s
- on every stripe Price object there must be a tier metadata
- Tier: an integer that quantifies the Price on a scalar axis. Example:
    - Monthly plan: tier = 1
    - Yearly plan: tier = 2