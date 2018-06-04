const ConfigService = angular.module('ConfigService', [])
    .service('config', ['$http', function ($http) {
        return $http.get('/config');
    }]);