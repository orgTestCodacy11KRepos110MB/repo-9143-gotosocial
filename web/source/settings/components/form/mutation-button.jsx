/*
	GoToSocial
	Copyright (C) 2021-2023 GoToSocial Authors admin@gotosocial.org

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

"use strict";

const React = require("react");

module.exports = function MutationButton({ label, result, disabled, ...inputProps }) {
	let iconClass = "";

	if (result.isLoading) {
		iconClass = "fa-spin fa-refresh";
	} else if (result.isSuccess) {
		iconClass = "fa-check fadeout";
	}

	return (<div>
		{result.error &&
			<section className="error">{result.error.status}: {result.error.data.error}</section>
		}
		<button type="submit" disabled={result.isLoading || disabled}	{...inputProps}>
			<i className={`fa fa-fw with-text ${iconClass}`} aria-hidden="true"></i>
			{result.isLoading
				? "Processing..."
				: label
			}
		</button>
	</div>
	);
};