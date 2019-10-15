pragma solidity ^0.5.8;

interface Upgradeable {
    function target() external view returns (address addr);
    function paused() external view returns (bool val);
    function owner() external view returns (address addr);
}
