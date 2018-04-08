# kits
Modular architecture for Go

This is an exercise in architecture design. I was looking for a modular approach
to building highly complicated software, like a SQL query engine, without a
monolithic tightly coupled code base. This is an attempt at a framework and API
that would support breaking up the code into many small packages/modules of code
(named `kits` to avoid confusion with Go's packages), where each implements a
simple API without knowledge or any hard links to the other kits. This allows
each such kit to be tested and used independently of the rest of the system,
with dependencies being natively and elegantly injected.
