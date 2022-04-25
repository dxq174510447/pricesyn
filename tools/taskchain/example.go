package taskchain

const Eg1 = `
name: ticket
version: 1
stages:
  - valid
  - ticketing
  - orderConfirm
  - voucherPrint
  - ticketSuccess
failure:
  - ticketFailure
`
