angular.module('BlocksApp').controller('RichlistController', function(
    $stateParams,
    $rootScope,
    $scope,
    $http
) {
    $rootScope.$state.current.data["pageSubTitle"] = $stateParams.hash;
    $scope.accounts = [];

    $http.get('/api/richlist?start=0&limit=100')
        .then(res => {
            $scope.accounts = res.data;
        });
})
