function NodeCtrl($scope, $http) {
  $scope.nodes = [];
  $scope.working = false;

  var logError = function(data, status) {
    console.log('code '+status+': '+data);
    $scope.working = false;
  };

  $scope.speak = function() {
    $scope.working = true;
    $http.get('/speak?s=' + $scope.saying).
      error(logError).
      success(function() {
        $scope.working = false;
        $scope.saying = '';
      });
  };

}
