angular.module('BlocksApp').controller('RichlistController', function(
    $stateParams,
    $rootScope,
    $scope,
    $http
) {
    $rootScope.$state.current.data["pageSubTitle"] = $stateParams.hash;
    $scope.richlist = {
        rankings: []
    };
    $scope.limit = 10;
    $scope.start = 0;
    $scope.isMoreDisabled = false;

    $scope.getMore = () => {
        $http.get(`/api/richlist?start=${$scope.start}&limit=${$scope.limit}`)
            .then(res => {
                $scope.richlist.rankings = $scope.richlist.rankings.concat(res.data.rankings);
                $scope.richlist.circulatingSupply = res.data.circulatingSupply;
                $scope.richlist.totalSupply = res.data.totalSupply;
                $scope.start += 10;

                if (res.data.rankings.length < $scope.limit) {
                    $scope.isMoreDisabled = true;
                }

                getSupplyOwned();
            });
    };

    $scope.getMore();

    function getSupplyOwned () {
        $scope.richlist.rankings = $scope.richlist.rankings.map(acct => {
            acct.supplyOwned = (acct.balance / $scope.richlist.circulatingSupply * 100).toFixed(2);

            return acct;
        });
    }
})
