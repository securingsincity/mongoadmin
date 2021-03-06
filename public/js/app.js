$(document).ready(function(){
  // $('.nav.nav-tabs a').click(function (e) {
  //   e.preventDefault();
  //   $(this).tab('show');
  // });
  $('ul.tabs').tabs();
});
angular.module('app', [])
.controller('MainCtrl', function($scope, $http) {
  $scope.findMore = {
    limit : 50,
    skip: 0
  };
  $http.get('/databases').success(function(response) {
    $scope.dbs = response;
  });

  $scope.setActive = function(db) {
    $scope.activeDb = db;
    $scope.activeCollection = null;
    $http.get('/databases/' + db.label + '/collections').success(function(response) {
      $scope.collections = response;
    });
  };


  var loadIndexes = function() {
    $http.get('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/indexes').success(function(response) {
      $scope.indexes = response;
    });
  };

  var loadFinds = function() {
    $http.post('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/find').success(function(response) {
      $scope.findResults = response;
    });
  };

  var getTotal = function() {
    $http.get('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/total').success(function(response) {
      $scope.total = response;
    });
  };
  $scope.$watch('activeCollection', function(newVal, oldVal) {
    if(newVal) {
      loadIndexes();
      loadFinds();
      getTotal();
    }
  });
  $scope.findMoreRecords = function() {
    var limit = $scope.findMore.limit ? $scope.findMore.limit : 50;
    var skip = $scope.findMore.skip ? $scope.findMore.skip : 0;
    var query = $scope.findMore.query ? JSON.parse($scope.findMore.query) : {};
    $http.post('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/find?limit=' + limit + '&skip=' + skip, {query: query}).success(function(response) {
      $scope.findResults = response;
    })
    .fail(function(error) {
      $scope.findResults = [error];
    });
  };
  $scope.insert = function() {
    var query = $scope.insert.doc ? JSON.parse($scope.insert.doc) : {};
    if (query === {}) {
      return false;
    }
    $http.post('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection , query).success(function(response) {
      $scope.insertResult = response;
    });
  };
  $scope.findByIdQuery = function() {
    var id = $scope.findById.id ? $scope.findById.id : '';
    $http.get('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/findById/'+ id)
      .success(function(response) {
        $scope.findByIdResults = response;
      })
      .fail(function(error) {
        $scope.findByIdResults = error;
      });
  };
  $scope.deleteByIdQuery = function(id) {
    if(confirm('Are you sure?')) {
      $http.delete('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/delete/'+ id)
        .success(function(response) {
          loadFinds();
        });
    }
  };
  $scope.addNewIndex = function() {
    var keys = $scope.newIndex.keys;
    var unique = $scope.newIndex.unique ? 'true' : 'false';
    var sparse = $scope.newIndex.sparse ? 'true' : 'false';
    if(keys) {
      $http.post('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/newIndex?keys=' + keys + '&unique=' + unique + '&sparse=' + sparse).success(function(response) {
        console.log(response);
        loadIndexes();
        $scope.newIndex = {};
      });
    }
  };


  $scope.removeIndex = function(idx) {
    console.log(idx);
    if(confirm('Are you sure?')) {
      $http.post('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/dropIndex?keys=' + idx.Key.join(',')).success(function(response) {
        console.log(response);
        loadIndexes();
      });
    }
  };
});