# Moving Google Cloud Firestore Data from Nested to Flat

This was a one-time use application but I wanted to store it for future reference.

## Old Model

I had originally nested data such that:

cats (collection) --> cats (documents) --> LitterTrips (collection) --> LitterTrips.
where cats had cat details like ID, Name, age, color, photoURL
and where
LitterTrips had details like time of trip, probability of trip, direction

## New Model

litterTrips (collection) --> litterTrips (documents)
where litterTrips have trip detail and CatID and CatName

cats (collection) still houses more detail about cats
catid, catname, color, age, photoURL

## Why change

The focus of the app(s) using the data is trips to the litterbox, not cats. The user wants to see the lates trips to the box regardless of cat.

We can still look for trips from a particular cat, but swapping the focus made the app much simpler.
