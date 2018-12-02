# Metrics
Metrics is a service to handle all kind of metrics (frontend, backend and system
metrics) in one service in order to make them available in different tools.

# Relevant make targets

* `make`: Build the binary and places it in the `bin` directory.
* `make update-vendors`: Obtain the external go dependencies and install them in
  the `vendor` directory.
* `make update-golden-files`: Update the testing golden files of an specific
  package. To target a specific package you need to set the `GOLDEN_PKG` variable.
  By default it updates the golden files for `github.com/radar-go/metrics`.
* `make tests`: Run the tests.
* `make lint`: Run the linting tests. If you want to generate a report file in
  the `reports` directory you need to set the variable `GENERATE_REPORT` to 1.
* `make coverage`: Run the tests code coverage for the project. If you want to
  generate a report file in the `reports` directory you need to set the variable
  `GENERATE_REPORT` to 1.
* `make start`: Starts the service locally in the port indicated in the
  configuration file.
* `make stop`: Stops the service running locally.
* `make restart`: Restarts the service running locally.

# License
metrics is licensed under the [GNU GPLv3](https://www.gnu.org/licenses/gpl.html).
You should have received a copy of the GNU General Public License along with
metrics. If not, see http://www.gnu.org/licenses/.

<p align="center">
<img src="https://www.gnu.org/graphics/gplv3-127x51.png" alt="GNU GPLv3">
</p>
