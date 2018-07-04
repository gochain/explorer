angular.module('BlocksApp').controller('RichlistController', function(
    $stateParams,
    $rootScope,
    $scope,
    $http
) {
    $rootScope.$state.current.data["pageSubTitle"] = $stateParams.hash;
    $scope.richlist = [];

    $http.get('/api/richlist?start=0&limit=100')
        .then(res => {
            $scope.richlist = res.data;
            $scope.richlist.rankings = $scope.richlist.rankings.map(acct => {
                acct.supplyOwned = supplyOwned(acct);

                return acct;
            })
        });

    function supplyOwned (account) {
        return (account.balance / $scope.richlist.circulatingSupply * 100).toFixed(2)
    }
})
