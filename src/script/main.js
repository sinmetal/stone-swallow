(function () {
    var app = angular.module('stone-swallow', ['ngResource']).
        config(function ($routeProvider) {
            $routeProvider.
                when('/', {controller: 'EntityListController', templateUrl: '/html/entity/list.html'})
        });

    app.directive('watchPath', ['$location', function ($location) {
        return function ($scope, $el, $attrs) {
            $scope.$on('$routeChangeSuccess', function () {
                var path = $location.path().split('/')[1];
                $el.toggleClass('active', path === $attrs.watchPath);
            });
        };
    }]);

    app.controller('EntityListController', ['$scope', '$resource', function ($scope, $resource) {
        $scope.search = function () {
            var entity = $resource("/entity?kind=" + $scope.kind);
                $scope.entities = entity.query(function () {
                    console.log("success entity query");
                    console.log($scope.entities);
                }, function () {
                    console.log("error entity query");
                });
        };
    }]);
})();