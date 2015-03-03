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
        $scope.limit = 1000;

        $scope.search = function () {
            var order = "";
            if ($scope.order) {
                order = $scope.order;
            }

            var limit = 0;
            if ($scope.limit) {
                limit = $scope.limit
            }
            var entity = $resource("/entity?kind=" + $scope.kind + "&order=" + order + "&limit=" + limit);
                $scope.entities = entity.query(function () {
                    console.log("success entity query");
                    console.log($scope.entities);
                }, function () {
                    console.log("error entity query");
                });
        };
    }]);
})();