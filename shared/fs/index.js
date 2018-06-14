// @flow
import * as I from 'immutable'
import * as React from 'react'
import * as Types from '../constants/types/fs'
import {globalStyles, globalMargins} from '../styles'
import {Box, Icon, List, ScrollView, Text} from '../common-adapters'
import FolderHeader from './header/container'
import SortBar from './sortbar/container'
import {Row, Placeholder} from './row'
import Footer from './footer/container'
import {isMobile} from '../constants/platform'
import ConnectedBanner from './banner/container'

type FolderProps = {
  items: Array<Types.Path>,
  isUserReset: boolean,
  resetParticipants: Array<string>,
  path: Types.Path,
  routePath: I.List<string>,
  progress: 'pending' | 'loaded',
}

class Files extends React.PureComponent<FolderProps> {
  _renderRow = (index, item) => (
    <Row key={Types.pathToString(item)} path={item} routePath={this.props.routePath} />
  )
  _renderRowPlaceholder = index => <Placeholder key={index} />

  render() {
    const content = this.props.isUserReset ? (
      <Box style={globalStyles.flexBoxColumn}>
        <Box style={resetContainerStyle}>
          <Icon type={isMobile ? 'icon-skull-64' : 'icon-skull-48'} />
          <Icon type="icon-access-denied-266" />
        </Box>
      </Box>
    ) : this.props.progress === 'pending' ? (
      <List items={['1', '2', '3']} renderItem={this._renderRowPlaceholder} />
    ) : this.props.items.length === 0 ? (
      <Box style={stylesEmptyContainer}>
        <Text type="BodySmall">This is an empty folder.</Text>
      </Box>
    ) : (
      <List items={this.props.items} renderItem={this._renderRow} />
    )
    return (
      <Box style={styleOuterContainer}>
        <Box style={stylesContainer}>
          <FolderHeader path={this.props.path} routePath={this.props.routePath} />
          <SortBar path={this.props.path} />
          {isMobile && this.props.resetParticipants.length > 0 ? (
            <ScrollView>
              <ConnectedBanner path={this.props.path} />
              <Box>
                {content}
              </Box>
            </ScrollView>
          ) : content
          }
          <Footer />
        </Box>
      </Box>
    )
  }
}

const styleOuterContainer = {
  height: '100%',
  position: 'relative',
}

const stylesContainer = {
  ...globalStyles.flexBoxColumn,
  ...globalStyles.fullHeight,
  flex: 1,
}

const stylesEmptyContainer = {
  ...globalStyles.flexBoxColumn,
  ...globalStyles.fullHeight,
  flex: 1,
  justifyContent: 'center',
  alignItems: 'center',
}

const resetContainerStyle = {
  ...globalStyles.flexBoxColumn,
  alignItems: 'center',
  flex: 1,
  justifyContent: 'center',
  marginTop: 2 * globalMargins.xlarge,
}

export default Files
