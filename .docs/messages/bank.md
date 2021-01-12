# x/bank

## MsgSend

```json
{
  "@type": "/cosmos.bank.v1beta1.MsgSend",
  "from_address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "to_address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
  "amount": [
    {
      "denom": "udaric",
      "amount": "1000"
    }
  ]
}
```

## MsgMultiSend

```json
{
  "@type": "/cosmos.bank.v1beta1.MsgMultiSend",
  "inputs": [
    {
      "address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
      "coins": [
        {
          "denom": "udaric",
          "amount": "1000"
        }
      ]
    }
  ],
  "outputs": [
    {
      "address": "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
      "coins": [
        {
          "denom": "udaric",
          "amount": "1000"
        }
      ]
    }
  ]
}
```
