nodeAffinity:
    required:
        nodeSelectorTerms:
        - matchFields: <>
          matchExpressions: <>
    preferred:
    -   weight:
        preference: 
            matchFields: <>
            matchExpressions: <>



<> means 
[]nodeSelectorRequirement:
    key, operator, []values   

nodeSelectorTerm:
    matchFields: <>
    matchExpressions: <>


In nodeAffinity,
multiple nodeSelectorTerms, under one type of affinity, are ORed.
multiple matchExpressions are ANDed.

----------------


"this Pod should (or, in the case of anti-affinity, should not) run in an X if that X is already running one or more Pods that meet rule Y"


* means
podAffinityTerm:
    labelSelector
    namespaces
    topologyKey
    namespaceSelector
    matchLabelKeys
    mismatchLabelKeys

podAffinity:
    required:
    - *
    preferred:
    - weight
      *

podAffinity:
    required:
    - *
    preferred:
    - weight
      *


multiple required terms :  all should be satisfied.
multiple preferred terms : calculated with summing the weights.


Important to note that :: 
multiple topologyConstraints are calculated individually &  inter-joined
affinity with topologyConstraints means filter out with affinity terms, then calculate topologyConstraints.