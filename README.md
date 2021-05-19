# namespace-admission-controller

作成したnamespaceを自動的にProjectに割り当てるadmission controller

## mutating admission-controller

- namespaceが作成された時にprojectを作成し、そのProjectIDをannotationとして付与するmutating-admission-controller

## validating admission-controller

- namespaceの作成がリクエストされた時に、そのユーザがProject作成の権限を持っているかどうかを確認し、持っていなければ弾くvalidating admission-controller
