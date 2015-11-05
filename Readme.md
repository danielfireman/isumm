# ISumm

I created this project based on a personal need of having a summarized view of my investiments. The idea is to have manually record monthly operations (i.e. balance, deposits or withdrawal) and properly generates a simple summary views of investments.

Main goals:

* Simple: no crazy JS or animations. As much stuff done server side as possible
* Super secure: users will have their own instance (app engine) of the whole app (include the complete datastore). The app will require google login and only one pre-configured user will have access to the instance
* Open source: just because :D 

## Main concepts

* The basic unit of time is month. Everything is (sooner or later) aggregated by month.
* Operation: something that happened to your investment account: Deposit, Withdrawal and Balance.
    * Withdrawal and Deposit: used to calculate month-by-month profit.
    * Balance: for the sake of simplification, only one is considered per month. If you added more than one balance operations, only the most-recent is used.
* Investment: yeah, the name says everything. It groups a set of operations.