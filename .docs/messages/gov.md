# x/gov

## MsgSubmitProposal

The `MsgSubsmitProposal` allows to submit a governance proposal.

There are various types of proposals that can be submitted. What distinguishes between one or another is the value of
the `content.@type` field. Following you can find the most common ones.

**Note**. Each chain can define its own proposal types.

### TextProposal

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgSubmitProposal",
  "content": {
    "@type": "/cosmos.gov.v1beta1.TextProposal",
    "title": "test software upgrade proposal",
    "description": "something about the proposal here"
  },
  "initial_deposit": [
    {
      "denom": "udaric",
      "amount": "20000000"
    }
  ],
  "proposer": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu"
}
```

### SoftwareUpgradeProposal

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgSubmitProposal",
  "content": {
    "@type": "/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal",
    "title": "test software upgrade proposal",
    "description": "something about the proposal here",
    "plan": {
      "name": "test",
      "time": "0001-01-01T00:00:00Z",
      "height": "50",
      "info": "",
      "upgraded_client_state": null
    }
  },
  "initial_deposit": [
    {
      "denom": "udaric",
      "amount": "20000000"
    }
  ],
  "proposer": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu"
}
```

### ParameterChangeProposal

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgSubmitProposal",
  "content": {
    "@type": "/cosmos.params.v1beta1.ParameterChangeProposal",
    "title": "Param change proposal",
    "description": "This is a param change proposal",
    "changes": [
      {
        "subspace": "staking",
        "key": "max_validators",
        "value": "200"
      }
    ]
  },
  "initial_deposit": [
    {
      "denom": "udaric",
      "amount": "20000000"
    }
  ],
  "proposer": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu"
}
```

### CommunityPoolSpendProposal

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgSubmitProposal",
  "content": {
    "@type": "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal",
    "title": "Community spend proposal",
    "description": "This is a community spend proposal",
    "recipient": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
    "amount": [
      {
        "denom": "udaric",
        "amount": "20000000"
      }
    ]
  },
  "initial_deposit": [
    {
      "denom": "udaric",
      "amount": "20000000"
    }
  ],
  "proposer": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu"
}
```

## MsgDeposit

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgDeposit",
  "proposal_id": "1",
  "depositor": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "amount": [
    {
      "denom": "udaric",
      "amount": "100"
    }
  ]
}
```

## MsgVote

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgVote",
  "proposal_id": "1",
  "voter": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "option": "VOTE_OPTION_YES"
}
```

