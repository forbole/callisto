# x/staking

## MsgCreateValidator

```json
{
  "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
  "description": {
    "moniker": "Validator moniker",
    "identity": "",
    "website": "",
    "security_contact": "",
    "details": ""
  },
  "commission": {
    "rate": "0.100000000000000000",
    "max_rate": "0.200000000000000000",
    "max_change_rate": "0.010000000000000000"
  },
  "min_self_delegation": "1",
  "delegator_address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "validator_address": "desmosvaloper13yp2fq3tslq6mmtq4628q38xzj75ethz8j43kw",
  "pubkey": {
    "@type": "/cosmos.crypto.ed25519.PubKey",
    "key": "1pk2pQfffJGLUqoOKQpHz1qnil0ymzYPEdSMufr1vTw="
  },
  "value": {
    "denom": "udaric",
    "amount": "1000000"
  }
}
```

## MsgDelegate

```json
{
  "@type": "/cosmos.staking.v1beta1.MsgDelegate",
  "delegator_address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "validator_address": "desmosvaloper13yp2fq3tslq6mmtq4628q38xzj75ethz8j43kw",
  "amount": {
    "denom": "udaric",
    "amount": "1000"
  }
}
```

## MsgEditValidator

```json
{
  "@type": "/cosmos.staking.v1beta1.MsgEditValidator",
  "description": {
    "moniker": "New moniker",
    "identity": "[do-not-modify]",
    "website": "[do-not-modify]",
    "security_contact": "[do-not-modify]",
    "details": "[do-not-modify]"
  },
  "validator_address": "desmosvaloper13yp2fq3tslq6mmtq4628q38xzj75ethz8j43kw",
  "commission_rate": null,
  "min_self_delegation": null
}
```

## MsgRedelegate

```json
{
  "@type": "/cosmos.staking.v1beta1.MsgBeginRedelegate",
  "delegator_address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "validator_src_address": "desmosvaloper13yp2fq3tslq6mmtq4628q38xzj75ethz8j43kw",
  "validator_dst_address": "desmosvaloper13yp2fq3tslq6mmtq4628q38xzj75ethz8j43kw",
  "amount": {
    "denom": "udaric",
    "amount": "1000"
  }
}
```

## MsgUndelegate

```json
{
  "@type": "/cosmos.staking.v1beta1.MsgUndelegate",
  "delegator_address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "validator_address": "desmosvaloper13yp2fq3tslq6mmtq4628q38xzj75ethz8j43kw",
  "amount": {
    "denom": "udaric",
    "amount": "1000"
  }
}
```