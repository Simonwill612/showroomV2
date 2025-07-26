// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract SmartTableControl {
    uint256 public targetHeight;
    bool public relayOn;
    mapping(string => string) public tableIdCommand;
    // ===================== EVENTS =====================
    event RaiseTable(string tableId);
    event LowerTable(string tableId);
    event StopTable(string tableId);
    event HeightSet(string tableId, uint256 height);
    event TableCommand(string tableId, string command);
    event EmitTableData(string tableId, string heightsCm);

    // ===================== TABLE COMMAND INTERFACE =====================
   function raiseTable(string memory tableId) public {
    emit RaiseTable(tableId);
    emit TableCommand(tableId, "UP");
    tableIdCommand[tableId] = "UP"; // 👈 thêm dòng này
}

function lowerTable(string memory tableId) public {
    emit LowerTable(tableId);
    emit TableCommand(tableId, "DOWN");
    tableIdCommand[tableId] = "DOWN";
} 

function stopTable(string memory tableId) public {
    emit StopTable(tableId);
    emit TableCommand(tableId, "STOP");
    tableIdCommand[tableId] = "STOP";
}

function setHeight(string memory tableId, uint256 height) public {
    targetHeight = height;
    string memory command = string(abi.encodePacked("HEIGHT=", uint2str(height)));
    emit HeightSet(tableId, height);
    emit TableCommand(tableId, command);
    tableIdCommand[tableId] = command;
}


    // ===================== TABLE DATA STORAGE =====================
    struct TableData {
        string tableId;
        string heightsCm;
    }

    mapping(string => TableData) public tableData;

    function tableHandleAI(string memory tableId, string memory heightsCm)
        public
    {
        tableData[tableId] = TableData({
            tableId: tableId,
            heightsCm: heightsCm
        });

        emit EmitTableData(tableId, heightsCm);
    }

    // ===================== UTILS =====================
    function uint2str(uint256 _i) internal pure returns (string memory str) {
        if (_i == 0) return "0";
        uint256 j = _i;
        uint256 length;
        while (j != 0) {
            length++;
            j /= 10;
        }
        bytes memory bstr = new bytes(length);
        uint256 k = length - 1;
        while (_i != 0) {
            bstr[k--] = bytes1(uint8(48 + (_i % 10)));
            _i /= 10;
        }
        str = string(bstr);
    }
}
