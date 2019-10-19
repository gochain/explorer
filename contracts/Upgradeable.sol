pragma solidity ^0.5.8;

interface Upgradeable {
    event Upgraded(address indexed target);
    event Paused();
    event Resumed();

    function target() external view returns (address addr);
    function paused() external view returns (bool val);
    function owner() external view returns (address addr);
}
