package taskchain

const Eg1 = `
name: ticket
version: 1
stage:
  - valid
  - ticketing
  - orderConfirm
  - voucherPrint
  - ticketSuccess
failure:
  - ticketFailure
`
