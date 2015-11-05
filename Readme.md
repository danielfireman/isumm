# ISumm

I created this project based on a personal need of having a summarized view of my investiments. The idea is to have a very simple manually record (monthly) balance, deposits or withdrawal and properly generates summary views of my investments.

It requires google login and the medium term goal is to make it simple so a person could just git clone this reposity, change two configuration files and reploy it to app engine. In this way users will have their own servers and datastores.

## Main concepts

* The basic unit of time is month. Everything is (sooner or later) aggregated by month.
* Operation: something that happened to your investment account: Deposit, Withdrawal and Balance.
    * Withdrawal and Deposit: used to calculate month-by-month profit.
    * Balance: for the sake of simplification, only one is considered per month. If you added more than one balance operations, only the most-recent is used.
* Investment: yeah, the name says everything. It groups a set of operations.