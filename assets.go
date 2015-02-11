package main

var HTML = `
<html>
	<head>
		<link rel="stylesheet" type="text/css" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css"/>
		<script src="https://ajax.googleapis.com/ajax/libs/angularjs/1.3.11/angular.min.js"></script>
		<script src="/app.js"></script>
	</head>
	<body ng-app="app">
		
		<div class="container-fluid" ng-controller="MainCtrl">
			<div class="row">
				<div class="col-sm-12">
					<h1>Manage DBs.</h1>
				</div>
			</div>
			<div class="row">
				<div class="col-sm-3">
					<div class="list-group">
						<a class="list-group-item" ng-repeat="db in dbs" ng-click="setActive(db)">
							<h4 class="list-group-item-heading">{{db.label}}</h4>
							
						</a>
					</div>
				</div>
				<div class="col-sm-9" ng-show="activeDb">
					<h2>{{activeDb.label}}</h2>
					<select ng-options="col for col in collections" ng-model="activeCollection" class="form-control">
						<option value="">Pick a collection.</option>
					</select>
					<hr />
					<div ng-show="activeCollection" class="panel panel-default">
						<div class="panel-heading">
							<h3 class="panel-title">{{activeCollection}}</h3>
							<h5>{{total | number}} total records</h5>
						</div>
						<div class="panel-body">
							<div class="row">
								<div class="col-sm-6">
									<h5>Indexes</h5>
									<ul class="list-group">
										<li class="list-group-item" ng-repeat="idx in indexes">
											<button class="btn btn-sm pull-right btn-danger" ng-click="removeIndex(idx)">X</button>
											<h4 class="list-group-item-heading">{{idx.Name}}</h4>
											<pre>{{idx.Key | json}}</pre>
											Unique: {{idx.Unique | json}}<br/>
											Sparse: {{idx.Sparse | json}}
										</li>
									</ul>
									
								</div>
								<div class="col-sm-6">
									<h6>Make an index</h6>
									<div class="form-group">
										<label>Keys</label>
										<input type="text" class="form-control input-sm" ng-model="newIndex.keys"/>
										<p class="help-block">(use "-" for descending, separate multiple keys with ",")</p>
									</div>
									<div class="checkbox">
										<label><input type="checkbox" ng-model="newIndex.unique"/> Unique</label>
									</div>
									<div class="checkbox">
										<label><input type="checkbox" ng-model="newIndex.sparse"/> Sparse</label>
									</div>
									
									<button ng-disabled="!newIndex.keys" class="btn btn-primary" ng-click="addNewIndex()">Add Index</button>
								</div>

							</div>
						</div>
					</div>

					
				</div>
			</div>
		</div>
	</body>

</html>
`

var JAVASCRIPT = `
angular.module('app', [])
.controller('MainCtrl', function($scope, $http) {
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

	var getTotal = function() {
		$http.get('/databases/' + $scope.activeDb.label + '/collections/' + $scope.activeCollection + '/total').success(function(response) {
			$scope.total = response;
		});
	}
	$scope.$watch('activeCollection', function(newVal, oldVal) {
		if(newVal) {
			loadIndexes();
			getTotal();
		}
	});

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
	}
});
`
