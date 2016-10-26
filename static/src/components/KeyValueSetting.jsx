import React from 'react'
import { Input, Button } from 'antd'
import { Box } from 'react-polymer-layout'
import { KVGet, KVPut, KVDelete } from './request'
import { message } from 'antd'

const KeyValueSetting = React.createClass({
    _getDone(result) {
        this.setState({ value: result.value })
    },

    _get(key) {
        KVGet(key || this.props.currentKey, this._getDone)
    },

    _updateDone(result) {
        message.info("update successfully.")
    },

    _update() {
        KVPut(this.props.currentKey, this.state.value, this._updateDone)
    },

    _deleteDone(result) {
        this.props.delete()
    },

    _delete() {
        KVDelete(this.props.currentKey, this._deleteDone)
    },

    getInitialState() {
        return { value: "" }
    },

    _fetch(key) {
        this.setState({ value: "" })
        this._get(key)
    },

    componentDidMount() {
        this._fetch()
    },

    componentWillReceiveProps(nextProps) {
        if (this.props.currentKey !== nextProps.currentKey) {
            this._fetch(nextProps.currentKey)
        }
    },

    render() {
        let mainColor = "#8ddafd"
        return (
            <Box vertical className="kv-editor" style={{ borderColor: mainColor }}>
                <div style={{ height: 20, backgroundColor: mainColor }}></div>
                <Box center style={{ height: 50, fontSize: 20, fontWeight: 500, borderBottom: "1px solid #ddd", paddingLeft: 5 }}>
                    {this.props.currentKey}
                </Box>
                <Box vertical style={{ padding: "10px 7px 0px 7px" }}>
                    <div style={{ width: "100%", paddingTop: 10 }}>
                        <Input type="textarea" rows={4} value={this.state.value} onChange={e => this.setState({ value: e.target.value })} />
                    </div>
                    <Box>
                        <Button type="primary" onClick={this._update} >Update</Button>
                        {
                            <Button type="ghost" onClick={this._delete} >Delete</Button>
                        }
                    </Box>
                </Box>
            </Box>
        )
    }
})

module.exports = KeyValueSetting